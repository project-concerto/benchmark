package org.example.e2eb.metric;

import org.HdrHistogram.Histogram;
import org.example.e2eb.Config;
import org.example.e2eb.emulator.Txns;

import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.atomic.AtomicInteger;

import static org.example.e2eb.utils.Utils.getSimpleClassName;


/**
 * Class to store the performance data
 * of an emulator
 */
public class Metric {

    /**
     * identify
     */
    private final long id;

    /**
     * warmUp start time, benchmark run and end time
     */
    public long warmStart = 0L;
    public long runStart = 0L;
    public long runEnd = 0L;
    public long coolDownEnd = 0L;

    /**
     * successful transaction latency
     */
    public final HashMap<String, Histogram> sLatencies = new HashMap<>();

    /**
     * successful transaction average and percentile latency
     */
    public final HashMap<String, Double> sAveLatency = new HashMap<>();
    public final HashMap<String, Long> sP50Latency = new HashMap<>();
    public final HashMap<String, Long> sP90Latency = new HashMap<>();
    public final HashMap<String, Long> sP99Latency = new HashMap<>();

    /**
     * successful transaction count and throughput
     */
    public final HashMap<String, Long> sCount = new HashMap<>();
    public final HashMap<String, Double> sThroughput = new HashMap<>();

    /**
     * failed transaction latency
     */
    public final HashMap<String, Histogram> fLatencies = new HashMap<>();

    /**
     * failed transaction average and percentile latency
     */
    public final HashMap<String, Double> fAveLatency = new HashMap<>();
    public final HashMap<String, Long> fP50Latency = new HashMap<>();
    public final HashMap<String, Long> fP90Latency = new HashMap<>();
    public final HashMap<String, Long> fP99Latency = new HashMap<>();

    /**
     * failed transaction count and throughput
     */
    public final HashMap<String, Long> fCount = new HashMap<>();
    public final HashMap<String, Double> fThroughput = new HashMap<>();

    /**
     * Constructor, initial performance data hashmap according types transaction
     */
    public Metric(int id) {
        this.id = id;
        Txns.getTxns().forEach(txnClass -> {
            sLatencies.put(getSimpleClassName(txnClass), new Histogram(3));
            sCount.put(getSimpleClassName(txnClass), 0L);
            sThroughput.put(getSimpleClassName(txnClass), 0.0);
            fLatencies.put(getSimpleClassName(txnClass), new Histogram(3));
            fCount.put(getSimpleClassName(txnClass), 0L);
            fThroughput.put(getSimpleClassName(txnClass), 0.0);
        });
    }

    /**
     * Calculate throughput and percentile of latency
     */
    public void summary() {
        Txns.getTxns().forEach(txnClass -> {
            calTput(txnClass, sThroughput, sCount);
            calLat(txnClass, sLatencies, sAveLatency, sP50Latency, sP90Latency, sP99Latency);
            calTput(txnClass, fThroughput, fCount);
            calLat(txnClass, fLatencies, fAveLatency, fP50Latency, fP90Latency, fP99Latency);
        });
        if (Config.getOptions().outPerEmulator) {
            System.out.print(String.join(", ",
                    String.valueOf(id), String.valueOf(runStart), String.valueOf(runEnd),
                    String.valueOf(runEnd - runStart), sCount.toString(), sThroughput.toString(),
                    sAveLatency.toString(), sP50Latency.toString(), sP90Latency.toString(), sP99Latency.toString(),
                    fCount.toString(), fThroughput.toString(), fAveLatency.toString(), fP50Latency.toString(),
                    fP90Latency.toString(), fP99Latency.toString()) + "\n");
        }
    }

    public void calLat(Class txnClass, HashMap<String, Histogram> latencies,
                       HashMap<String, Double> aveLatency, HashMap<String, Long> p50Latency,
                       HashMap<String, Long> p90Latency, HashMap<String, Long> p99Latency) {
        aveLatency.put(getSimpleClassName(txnClass), latencies.get(getSimpleClassName(txnClass)).getMean());
        p50Latency.put(getSimpleClassName(txnClass), latencies.get(getSimpleClassName(txnClass)).getValueAtPercentile(50));
        p90Latency.put(getSimpleClassName(txnClass), latencies.get(getSimpleClassName(txnClass)).getValueAtPercentile(90));
        p99Latency.put(getSimpleClassName(txnClass), latencies.get(getSimpleClassName(txnClass)).getValueAtPercentile(99));

    }

    public void calTput(Class txnClass, HashMap<String, Double> throughput, HashMap<String, Long> count){
        throughput.put(getSimpleClassName(txnClass), (count.get(getSimpleClassName(txnClass)) * 1000.0 / (runEnd - runStart)));
    }
}
