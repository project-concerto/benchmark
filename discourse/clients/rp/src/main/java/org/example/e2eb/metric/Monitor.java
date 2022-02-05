package org.example.e2eb.metric;

import org.example.e2eb.Config;
import org.example.e2eb.emulator.BaseEmulator;
import org.example.e2eb.emulator.Txns;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.concurrent.atomic.AtomicInteger;

import static org.example.e2eb.utils.Utils.getSimpleClassName;


/**
 * Class used to monitor the over all performance of all clients
 * and give a summary after benchmarking
 */
public class Monitor extends Thread {
    protected Logger logger = LoggerFactory.getLogger(BaseEmulator.class);

    /**
     * Emulators to be monitored
     */
    private List<BaseEmulator> emulators;

    /**
     * Global txn count for 0->warm up, 1->benchmark and 2->cool down
     */
    public static List<HashMap<String, AtomicInteger>> globalCnt =
            Arrays.asList(new HashMap<>(), new HashMap<>(), new HashMap<>());

    /**
     * interval time to output overall performance
     */
    private int interval;

    /**
     * Constructor
     *
     * @param emulators the emulators to be monitored
     */
    public Monitor(List<BaseEmulator> emulators) {
        this.emulators = emulators;
        this.metric = new Metric(0);
        this.interval = Config.getOptions().interval;
        Txns.getTxns().forEach(txnClass -> {
            globalCnt.get(0).put(getSimpleClassName(txnClass), new AtomicInteger(0));
            globalCnt.get(1).put(getSimpleClassName(txnClass), new AtomicInteger(0));
            globalCnt.get(2).put(getSimpleClassName(txnClass), new AtomicInteger(0));
        });
    }

    /**
     * Global metric
     */
    public Metric metric;

    /**
     * output the overall performance of benchmark
     */
    public void summary() {
        long minWarmUpStart = System.currentTimeMillis(), minRunStart = System.currentTimeMillis(),
                    maxRunEnd = 0, maxCoolDownEnd = 0;
        if(Config.getOptions().outPerEmulator){
            System.out.println("-------------------------Per Thread Performance-------------------------");
            System.out.println("Thread Id, Start Time, End Time, Duration, sCnt, sTput, sAveLat, sP50Lat, sP90Lat, " +
                    "sP99Lat, fCnt, fTput, fAveLat, fP50Lat, fP90Lat, fP99Lat");
        }
        for (BaseEmulator emulator : emulators) {
            emulator.getMetric().summary();
            Txns.getTxns().forEach(txnClass -> {
                metric.sCount.put(getSimpleClassName(txnClass),
                        metric.sCount.get(getSimpleClassName(txnClass)) + emulator.getMetric().sCount.get(getSimpleClassName(txnClass)));
                metric.sThroughput.put(getSimpleClassName(txnClass),
                        metric.sThroughput.get(getSimpleClassName(txnClass)) + emulator.getMetric().sThroughput.get(getSimpleClassName(txnClass)));
                metric.sLatencies.get(getSimpleClassName(txnClass)).add(emulator.getMetric().sLatencies.get(getSimpleClassName(txnClass)));
                metric.calLat(txnClass, metric.sLatencies, metric.sAveLatency, metric.sP50Latency, metric.sP90Latency, metric.sP99Latency);
                metric.fCount.put(getSimpleClassName(txnClass),
                        metric.fCount.get(getSimpleClassName(txnClass)) + emulator.getMetric().fCount.get(getSimpleClassName(txnClass)));
                metric.fThroughput.put(getSimpleClassName(txnClass),
                        metric.fThroughput.get(getSimpleClassName(txnClass)) + emulator.getMetric().fThroughput.get(getSimpleClassName(txnClass)));
                metric.fLatencies.get(getSimpleClassName(txnClass)).add(emulator.getMetric().fLatencies.get(getSimpleClassName(txnClass)));
                metric.calLat(txnClass, metric.fLatencies, metric.fAveLatency, metric.fP50Latency, metric.fP90Latency, metric.fP99Latency);

            });
            minWarmUpStart = Long.min(minWarmUpStart, emulator.getMetric().warmStart);
            minRunStart = Long.min(minRunStart, emulator.getMetric().runStart);
            maxRunEnd = Long.max(maxRunEnd, emulator.getMetric().runEnd);
            maxCoolDownEnd = Long.max(maxCoolDownEnd, emulator.getMetric().coolDownEnd);
        }
        if(!(minWarmUpStart <= minRunStart && minRunStart <= maxRunEnd && maxRunEnd <= maxCoolDownEnd)){
            logger.error("Some emulators' benchmark time exceed others' lifetime, which means contention is not fixed! " +
                    "{}, {}, {}, {}", minWarmUpStart, minRunStart, maxRunEnd, maxCoolDownEnd);
        }
        System.out.printf("-------------------------Benchmark Result Summary-------------------------\n" +
                        "Emulators Number: %d, Warm up: %d, Benchmark: %d, Cool down: %d\n" +
                        "sCount: %s, sThroughput: %s\n" +
                        "sAveLat: %s, sP50Lat: %s, sP90Lat: %s, sP99Lat: %s\n" +
                        "fCount: %s, fThroughput: %s\n" +
                        "fAveLat: %s, fP50Lat: %s, fP90Lat: %s, fP99Lat: %s\n",
                emulators.size(),
                Config.getOptions().warmUp, Config.getOptions().benchmark, Config.getOptions().coolDown,
                metric.sCount, metric.sThroughput,
                metric.sAveLatency, metric.sP50Latency, metric.sP90Latency, metric.sP99Latency,
                metric.fCount, metric.fThroughput,
                metric.fAveLatency, metric.fP50Latency, metric.fP90Latency, metric.fP99Latency);
        System.out.printf("-------------------------Benchmark Result CSV-------------------------\n" +
                        "Thread Num, Warm up, Benchmark, Cool down, sCnt, sTput, sAveLat, sP50Lat, sP90Lat, sP99Lat, " +
                        "fCnt, fTput, fAveLat, fP50Lat, fP90Lat, fP99Lat\n" +
                        "%d, %d, %d, %d, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s\n",
                emulators.size(),
                Config.getOptions().warmUp, Config.getOptions().benchmark, Config.getOptions().coolDown,
                metric.sCount, metric.sThroughput,
                metric.sAveLatency, metric.sP50Latency, metric.sP90Latency, metric.sP99Latency,
                metric.fCount, metric.fThroughput,
                metric.fAveLatency, metric.fP50Latency, metric.fP90Latency, metric.fP99Latency);
    }

    /**
     * Output overall performance every @interval
     */
    @Override
    public void run() {
        long warmUp = (long) Config.getOptions().warmUp;
        long benchmark = (long) Config.getOptions().benchmark;
        long coolDown = (long) Config.getOptions().coolDown;
        HashMap<String, Integer> last = new HashMap<>();
        HashMap<String, Integer> delta = new HashMap<>();
        Txns.getTxns().forEach(txnClass -> {
            last.put(getSimpleClassName(txnClass), 0);
            delta.put(getSimpleClassName(txnClass), 0);
        });
        System.out.printf("-------------------------Transaction Per %d s-------------------------\n", interval);
        // Warm up
        long start = System.currentTimeMillis();
        while(processNotStop(warmUp, start, 0)){
            try {
                Thread.sleep(interval * 1000L);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
            calDelta(last, delta, globalCnt.get(0));
            System.out.println("Warm up: " + delta);
        }
        Txns.getTxns().forEach(txnClass -> {
            last.put(getSimpleClassName(txnClass), 0);
        });

        start = System.currentTimeMillis();
        while (processNotStop(benchmark, start, 1)) {
            try {
                Thread.sleep(interval * 1000L);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
            calDelta(last, delta, globalCnt.get(1));
            System.out.println("Benchmark: " + delta);
        }
        Txns.getTxns().forEach(txnClass -> {
            last.put(getSimpleClassName(txnClass), 0);
        });

        start = System.currentTimeMillis();
        while(processNotStop(coolDown, start, 2)){
            try {
                Thread.sleep(interval * 1000L);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
            calDelta(last, delta, globalCnt.get(2));
            System.out.println("Cool down: " + delta);
        }
    }

    /**
     * Calculate the sum of transaction count in a hashmap
     * @param globalCnt global transaction count hashmap
     * @return the sum of different transaction count
     */
    public static int cntSum(HashMap<String, AtomicInteger> globalCnt) {
        return Txns.getTxns().stream().
                mapToInt(txnClass -> globalCnt.get(getSimpleClassName(txnClass)).get()).sum();
    }

    /**
     * Determine whether a process should stop
     * @param required How long should this process keeping, a time (useCnt = false) or count (useCnt = true)
     * @param processStart The start time of current process
     * @param type Process type, 0->warm up, 1->benchmark, 2->cool down
     * @return return false if current stop should stop, otherwise false
     */
    public static boolean processNotStop(Long required, Long processStart, Integer type) {
        if (Config.getOptions().useCnt) {
            return cntSum(globalCnt.get(type)) <= required;
        } else {
            return System.currentTimeMillis() - processStart <= required * 1000;
        }
    }

    /**
     * Calculate the cnt delta between an interval
     * @param last count before the interval
     * @param delta count difference
     * @param globalCnt current global count
     */
    public void calDelta(HashMap<String, Integer> last, HashMap<String, Integer> delta, HashMap<String, AtomicInteger> globalCnt) {
        Txns.getTxns().forEach(txnClass -> {
            String key = getSimpleClassName(txnClass);
            int curCnt = globalCnt.get(key).get();
            delta.put(key, curCnt - last.get(key));
            last.put(key, curCnt);
        });
    }
}
