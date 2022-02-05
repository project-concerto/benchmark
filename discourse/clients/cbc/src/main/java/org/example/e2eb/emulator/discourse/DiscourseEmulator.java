package org.example.e2eb.emulator.discourse;

import okhttp3.Headers;
import org.example.e2eb.emulator.BaseEmulator;
import org.example.e2eb.emulator.discourse.txns.CreatePostTxn;
import org.example.e2eb.emulator.discourse.txns.ToggleAnswerTxn;

public class DiscourseEmulator extends BaseEmulator {
    public final static String apiKey = "780b73d30f5057bd5ad9263cd25461ae162182b7a34aad55dc9eb549c14fb975";

    /**
     * construct an emulator, execute preTxn and choose a nextTxn to be executed
     *
     * @param eId the Id of current emulator
     */
    public DiscourseEmulator(int eId) {
        super(eId);
        this.isKeying = false;
        this.isThinking = false;
    }

    @Override
    public void preTxn() {

    }

    @Override
    public void nextTxn() {
        if(eId % 2 == 1){
            currentTxn = ToggleAnswerTxn.class;
        }
        else{
            currentTxn = CreatePostTxn.class;
        }
    }

    @Override
    public long getThinkTime() {
        return 0;
    }

    @Override
    public long getKeyingTime() {
        return 0;
    }
}
