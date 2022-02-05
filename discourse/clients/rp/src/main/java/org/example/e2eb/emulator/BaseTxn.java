package org.example.e2eb.emulator;

import org.example.e2eb.emulator.BaseEmulator;

/**
 * Base transaction class, a simple example
 * do not implement this class, because doTxn method should be static,
 * however java doesn't support abstract static method
 */
public abstract class BaseTxn {

    /**
     * do a transaction
     * @param emulator emulator that try to do this transaction
     * @return transaction execution result, true for successful and false for failed
     */
    public abstract boolean doTxn(BaseEmulator emulator);

}
