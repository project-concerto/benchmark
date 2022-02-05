package org.example.e2eb.emulator.broadleaf;

import org.example.e2eb.emulator.BaseEmulator;
import org.example.e2eb.emulator.broadleaf.txns.BLCTxn;
import org.example.e2eb.utils.Panic;

import java.io.IOException;

public class BLCEmulator extends BaseEmulator {

    /**
     * construct an emulator, execute preTxn and choose a nextTxn to be executed
     *
     * @param eId
     */
    public BLCEmulator(int eId) {
        super(eId);
        this.isKeying = false;
        this.isThinking = false;
    }

    @Override
    public void preTxn() {
        try {
            BLCTxn.login(okHttpClient, eId + "@qq.com", "zxd123");
        } catch (IOException e){
            Panic.quit("Login failed for user " + eId);
        }
    }

    @Override
    public void nextTxn() {
        currentTxn = BLCTxn.class;
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
