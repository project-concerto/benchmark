package org.example.e2eb.emulator.repair.txns;

import okhttp3.Headers;
import okhttp3.Response;
import org.example.e2eb.Config;
import org.example.e2eb.emulator.BaseEmulator;
import org.example.e2eb.emulator.repair.RepairEmulator;
import org.example.e2eb.emulator.repair.postdata.EditPostData;
import org.example.e2eb.emulator.request.RequestUtils;
import org.example.e2eb.utils.Utils;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

public class EditPostTxn {
    private final static String createPostUrl = Config.getOptions().protocol + Config.getOptions().host + "/posts/";
    private final static List<Integer> ids = new ArrayList<>();
    static {
        for(int i=1; i<=256; i++){
            ids.add(i);
        }
    }
    public static boolean doTxn(BaseEmulator emulator) {
        int eId = 0;
        try {
            synchronized (ids) {
                int idx = Utils.randomInt(0, ids.size() - 1);
                eId = ids.get(idx);
                ids.remove(idx);
            }
            Headers headers = new Headers.Builder()
                    .add("Api-Key", RepairEmulator.apiKey)
                    .add("Api-Username", eId + "qqcom")
                    .build();
            int postId = 268 + eId;
            try {
                try (Response response = RequestUtils.sendPutRequest(emulator.getOkHttpClient(), createPostUrl + postId,
                        new EditPostData(eId), headers)) {
                    if (response.code() != 200) {
                        return false;
                    }
                }
            } catch (IOException ignored) {
                return false;
            }
            return true;
        }
        finally {
            synchronized (ids) {
                ids.add(eId);
            }
        }
    }
}
