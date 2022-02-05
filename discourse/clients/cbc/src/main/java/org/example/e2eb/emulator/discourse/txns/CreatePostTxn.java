package org.example.e2eb.emulator.discourse.txns;

import okhttp3.Headers;
import okhttp3.Response;
import org.example.e2eb.Config;
import org.example.e2eb.emulator.BaseEmulator;
import org.example.e2eb.emulator.discourse.DiscourseEmulator;
import org.example.e2eb.emulator.discourse.postdata.CreatePostData;
import org.example.e2eb.emulator.request.RequestUtils;

import java.io.IOException;

public class CreatePostTxn {
    private final static String createPostUrl = Config.getOptions().protocol + Config.getOptions().host + "/posts.json";

    public static boolean doTxn(BaseEmulator emulator) {
        Headers headers = new Headers.Builder()
                .add("Api-Key", DiscourseEmulator.apiKey)
                .add("Api-Username", emulator.geteId() + "qqcom")
                .build();
        // use commented code for no contention
        // int topicId = 9 + emulator.geteId()/2 + 128;
        int topicId = 9 + emulator.geteId()/2;
        try {
            try(Response response = RequestUtils.sendPostRequest(emulator.getOkHttpClient(), createPostUrl,
                    new CreatePostData(topicId), headers)){
                if (response.code() != 200) {
                    return false;
                }
            }
        } catch (IOException ignored) {
            return false;
        }
        return true;
    }
}
