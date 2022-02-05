package org.example.e2eb.emulator.discourse.postdata;

import okhttp3.FormBody;
import org.example.e2eb.emulator.request.BasePostData;
import org.example.e2eb.utils.Utils;

public class CreatePostData extends BasePostData {
    public CreatePostData(int topic_id) {
        data = new FormBody.Builder()
                .add("raw", Utils.randomSentence(50))
                .add("topic_id", String.valueOf(topic_id))
                .build();
    }
}
