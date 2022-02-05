package org.example.e2eb.emulator.discourse.postdata;

import okhttp3.FormBody;
import org.example.e2eb.emulator.request.BasePostData;

public class ToggleAnswerData extends BasePostData {
    public ToggleAnswerData(int post_id){
        data = new FormBody.Builder()
                .add("id", String.valueOf(post_id))
                .build();
    }
}
