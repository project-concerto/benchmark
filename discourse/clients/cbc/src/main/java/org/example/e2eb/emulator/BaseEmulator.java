package org.example.e2eb.emulator;

import okhttp3.Cookie;
import okhttp3.CookieJar;
import okhttp3.HttpUrl;
import okhttp3.OkHttpClient;
import org.example.e2eb.Config;
import org.example.e2eb.emulator.request.RequestUtils;
import org.example.e2eb.metric.Metric;
import org.example.e2eb.metric.Monitor;
import org.example.e2eb.utils.Panic;
import org.jetbrains.annotations.NotNull;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Method;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.concurrent.TimeUnit;

import static org.example.e2eb.utils.Utils.getSimpleClassName;

/**
 * Base Class of emulators, extend this class and
 * implement @preTxn and @nextTxn according to requirements
 */
public abstract class BaseEmulator extends Thread {

    protected Logger logger = LoggerFactory.getLogger(BaseEmulator.class);


    /**
     * Identity of the emulator
     */
    protected int eId = 0;

    /**
     * metric of this emulator
     */
    protected Metric metric = null;

    /**
     * time for warmUp, benchmark and coolDown
     */
    protected long warmUp = 0L;
    protected long benchmark = 0L;
    protected long coolDown = 0L;

    /**
     * current transaction class to be executed
     */
    protected Class currentTxn = null;

    /**
     * Whether sleep to simulate keying
     */
    protected boolean isKeying = false;

    /**
     * Whether sleep to simulate thinking
     */
    protected boolean isThinking = false;

    /**
     * http client for current thread
     */
    protected OkHttpClient okHttpClient;

    /**
     * http client cookie
     */
    protected final HashMap<String, List<Cookie>> cookies = new HashMap<>();

    /**
     * Whether use global to control the start and stop of benchmark
     */
    private final boolean useCnt;

    /**
     * execution result of last transaction
     */
    protected boolean ok = true;

    /**
     * construct an emulator, execute preTxn and choose a nextTxn to be executed
     *
     * @param eId the Id of current emulator
     */
    public BaseEmulator(int eId) {
        this.eId = eId;
        this.metric = new Metric(eId);
        this.warmUp = (long) Config.getOptions().warmUp;
        this.benchmark = (long) Config.getOptions().benchmark;
        this.coolDown = (long) Config.getOptions().coolDown;
        this.useCnt = Config.getOptions().useCnt;
        okHttpClient = RequestUtils.getInsecureOkHttpClientBuilder().cookieJar(new CookieJar() {
            @Override
            public void saveFromResponse(@NotNull HttpUrl httpUrl, @NotNull List<Cookie> list) {
                cookies.put(httpUrl.host(), list);
            }

            @NotNull
            @Override
            public List<Cookie> loadForRequest(@NotNull HttpUrl httpUrl) {
                List<Cookie> cookie = cookies.get(httpUrl.host());
                return cookie != null ? cookie : new ArrayList<>();
            }
        }).connectTimeout(30, TimeUnit.SECONDS)
                .readTimeout(30, TimeUnit.SECONDS)
                .writeTimeout(30, TimeUnit.SECONDS).build();
    }

    /**
     * execute some transaction before benchmark, such as register or login
     */
    public abstract void preTxn();

    /**
     * Choose a transaction to be executed through changing @currentTxn,
     * implement according to requirement
     */
    public abstract void nextTxn();


    /**
     * generate a time interval for simulate user think
     *
     * @return time to think
     */
    public abstract long getThinkTime();

    /**
     * generate a time interval for simulate user keying
     *
     * @return time to keying
     */
    public abstract long getKeyingTime();

    /**
     * invoke current transaction according to its class
     *
     * @return result of transaction execution, true for successful and false for failed
     */
    private boolean invokeTxnByClass() {
        boolean ret = false;
        try {
            Method method = currentTxn.getMethod("doTxn", BaseEmulator.class);
            ret = (Boolean) method.invoke(null, this);
        } catch (Exception e) {
            Panic.quit(e.toString());
        }
        return ret;
    }


    /**
     * Begin a benchmark
     */
    @Override
    public void run() {
        preTxn();
        nextTxn();
        metric.warmStart = System.currentTimeMillis();
        // Warmup first
        while (Monitor.processNotStop(warmUp, metric.warmStart, 0)) {
            ok = invokeTxnByClass();
            if (ok) {
                Monitor.globalCnt.get(0).get(getSimpleClassName(currentTxn)).addAndGet(1);
            }
            nextTxn();
        }

        metric.runStart = System.currentTimeMillis();
        String txnName;
        // Begin benchmark
        while (Monitor.processNotStop(benchmark, metric.runStart, 1)) {
            // Keying for current transaction
            if (isKeying) {
                try {
                    Thread.sleep(getKeyingTime());
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }

            // Execute current transaction
            long txnStart = System.nanoTime();
            ok = invokeTxnByClass();
            // Record result in metric
            txnName = getSimpleClassName(currentTxn);
            if (ok) {
                metric.sLatencies.get(txnName).recordValue(System.nanoTime() - txnStart);
                metric.sCount.put(txnName, metric.sCount.get(getSimpleClassName(currentTxn)) + 1);
                Monitor.globalCnt.get(1).get(txnName).addAndGet(1);
            } else {
                metric.fLatencies.get(txnName).recordValue(System.nanoTime() - txnStart);
                metric.fCount.put(txnName, metric.fCount.get(txnName) + 1);
            }

            if (isThinking) {
                // User think
                try {
                    Thread.sleep(getThinkTime());
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }

            nextTxn();
        }
        metric.runEnd = System.currentTimeMillis();

        while (Monitor.processNotStop(coolDown, metric.runEnd, 2)) {
            ok = invokeTxnByClass();
            if (ok) {
                Monitor.globalCnt.get(2).get(getSimpleClassName(currentTxn)).addAndGet(1);
            }
            nextTxn();
        }
        metric.coolDownEnd = System.currentTimeMillis();
    }

    public Metric getMetric() {
        return metric;
    }

    public OkHttpClient getOkHttpClient() {
        return okHttpClient;
    }

    public int geteId() {
        return eId;
    }

    public boolean getOk() {
        return this.ok;
    }
}
