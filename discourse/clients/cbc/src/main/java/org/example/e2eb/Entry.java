package org.example.e2eb;

import org.example.e2eb.emulator.BaseEmulator;
import org.example.e2eb.emulator.Txns;
import org.example.e2eb.emulator.discourse.DiscourseEmulator;
import org.example.e2eb.metric.Monitor;
import org.example.e2eb.utils.Panic;
import org.kohsuke.args4j.CmdLineParser;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.List;

public class Entry {
    private static Logger logger = LoggerFactory.getLogger(Entry.class);

    public static void main(String[] args) throws InterruptedException {
        CmdLineParser parser = new CmdLineParser(Config.getOptions());
        if (args.length != 0) {
            try {
                parser.parseArgument(args);
            } catch (Exception e) {
                logger.error(e + "Cannot parse arguments");
            }
        }

        initialTxnClass();

        int emulatorNum = Config.getOptions().emulators;
        List<BaseEmulator> emulators = new ArrayList<>();

        for(int i=1; i<=emulatorNum; i++){
            emulators.add(new DiscourseEmulator(i));
        }

        Monitor monitor = new Monitor(emulators);
        monitor.start();
        for(BaseEmulator em : emulators){
            em.start();
        }
        for(BaseEmulator em : emulators){
            em.join();
        }
        monitor.join();
        monitor.summary();
    }

    public static void initialTxnClass(){
        List<String> classNames = Config.getOptions().txns;
        List<String> weights = Config.getOptions().weight;
        List<Class> classes = new ArrayList<>();
        List<Double> numWeight = new ArrayList<>();
        for(int i=0; i<classNames.size(); i++){
            try {
                classes.add(Class.forName("org.example.e2eb.emulator." + classNames.get(i)));
            } catch (ClassNotFoundException e) {
                Panic.quit("Class not found " + e);
            }
            numWeight.add(Double.valueOf(weights.get(i)));
        }
        Txns.initial(classes.toArray(new Class[0]), numWeight.toArray(new Double[0]));
    }
}
