package org.example.e2eb.emulator.repair;

import org.example.e2eb.emulator.BaseEmulator;
import org.example.e2eb.emulator.repair.txns.EditPostTxn;

public class RepairEmulator extends BaseEmulator {
    public final static String apiKey = "b323b67f0ee419bb696133cf114db161f37f29f296583f24a411b45344ebecc7";

    /**
     * construct an emulator, execute preTxn and choose a nextTxn to be executed
     *
     * @param eId the Id of current emulator
     */
    public RepairEmulator(int eId) {
        super(eId);
        this.isKeying = false;
        this.isThinking = false;
    }

    @Override
    public void preTxn() {

    }

    @Override
    public void nextTxn() {
        currentTxn = EditPostTxn.class;
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
