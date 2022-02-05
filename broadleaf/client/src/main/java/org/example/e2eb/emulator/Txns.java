package org.example.e2eb.emulator;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.Random;

/**
 * Class to store the information of transaction types
 */
public class Txns {

    private static final Logger logger = LoggerFactory.getLogger(Txns.class);

    /**
     * transaction types
     */
    private static final List<Class> txns = new ArrayList<>();

    /**
     * weight to be chosen of each transaction type
     */
    private static final List<Double> weight = new ArrayList<>();

    /**
     * low bound helping to choose a type according to weight
     */
    private static final List<Double> low = new ArrayList<>();

    /**
     * up bound helping to choose a type according to weight
     */
    private static final List<Double> up = new ArrayList<>();


    /**
     * initial transaction types and weight, and calculate corresponding
     * up and low bound
     * @param txnList transaction type list
     * @param weightList transaction weight list
     */
    public static void initial(Class[] txnList, Double[] weightList){
        txns.addAll(Arrays.asList(txnList));
        weight.addAll(Arrays.asList(weightList));
        low.addAll(weight);
        up.addAll(weight);
        int len = txns.size();
        low.set(0, 0.0);
        up.set(len - 1, 1.0);
        for(int i=1; i<len; i++){
            double sum = 0.0;
            for(int j=0; j<i; j++){
                sum += weight.get(j);
            }
            low.set(i, sum);
            up.set(i - 1, sum);
        }
    }

    /**
     * select a transaction according to weight
     * @return transaction class
     */
    public static Class randomTxnByWeight(){
        Random random = new Random();
        double r = random.nextDouble();
        for(int i=0; i<5; i++){
            if(r >= low.get(i) && r < up.get(i)){
                return txns.get(i);
            }
        }
        logger.error("r {} can not find corresponding txns between low {} and up {}", r, low, up);
        return BaseTxn.class;
    }

    public static List<Class> getTxns() {
        return txns;
    }
}
