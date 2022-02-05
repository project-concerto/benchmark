package org.example.e2eb;

import org.kohsuke.args4j.Option;
import org.kohsuke.args4j.spi.StringArrayOptionHandler;

import java.util.Arrays;
import java.util.List;

public class Config {
    private static final CmdOption options = new CmdOption();

    public static class CmdOption {
        @Option(name = "--protocol", usage = "protocol used, http or https")
        public String protocol = "http://";

        @Option(name = "--host", usage = "host of server")
        public String host = "localhost:3000";

        @Option(name = "--emulators", usage = "number of emulators")
        public Integer emulators = 2;

        @Option(name = "--benchmark", usage = "time for running the benchmark")
        public Integer benchmark = 30;

        @Option(name = "--warmUp", usage = "time for warm up")
        public Integer warmUp = 0;

        @Option(name = "--coolDown", usage = "time for cool down")
        public Integer coolDown = 0;

        @Option(name = "--interval", usage = "interval time for result print")
        public Integer interval = 5;

        @Option(name = "--debug", usage = "whether output debug log")
        public boolean debug = false;

        @Option(name = "--outPerEmulator", usage = "Output result per emulator")
        public boolean outPerEmulator = false;

        @Option(name = "--txns", handler = StringArrayOptionHandler.class)
        public List<String> txns = Arrays.asList("discourse.txns.CreatePostTxn", "discourse.txns.ToggleAnswerTxn");

        @Option(name = "--weight", handler = StringArrayOptionHandler.class)
        public List<String> weight = Arrays.asList("0.5", "0.5");

        @Option(name = "--useCnt", usage = "whether use count as termination time")
        public boolean useCnt = false;

        @Option(name = "--exponent", usage = "exponent of zipf dist.")
        public double exponent = 0.5;
    }

    public static CmdOption getOptions() {
        return options;
    }
}

