package org.example.e2eb.emulator.discourse.txns;

import okhttp3.Headers;
import okhttp3.Response;
import org.example.e2eb.Config;
import org.example.e2eb.emulator.BaseEmulator;
import org.example.e2eb.emulator.discourse.DiscourseEmulator;
import org.example.e2eb.emulator.discourse.postdata.ToggleAnswerData;
import org.example.e2eb.emulator.request.RequestUtils;

import java.io.IOException;

public class ToggleAnswerTxn {
    private final static String setAnswerUrl = Config.getOptions().protocol +
            Config.getOptions().host + "/solution/accept";

    private final static ThreadLocal<Boolean> isFirst = new ThreadLocal<>();

    public static boolean doTxn(BaseEmulator emulator) {
        if (isFirst.get() == null) {
            isFirst.set(true);
        }
        Headers headers = new Headers.Builder()
                .add("Api-Key", DiscourseEmulator.apiKey)
                .add("Api-Username", emulator.geteId() + "qqcom")
                .build();

        int postId = 268 + emulator.geteId() + (isFirst.get() ? 0 : 1);
        boolean ret = true;
        try {
            try (Response response = RequestUtils.sendPostRequest(emulator.getOkHttpClient(), setAnswerUrl,
                    new ToggleAnswerData(postId), headers)) {
                if (response.code() != 200) {
                    ret = false;
                }
            }
        } catch (IOException ignored) {
            ret = false;
        }
        if ((ret && isFirst.get()) || (!ret && !isFirst.get())) {
            isFirst.set(false);
        } else {
            isFirst.set(true);
        }
        return ret;
    }
}
