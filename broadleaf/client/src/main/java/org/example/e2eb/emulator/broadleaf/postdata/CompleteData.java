package org.example.e2eb.emulator.broadleaf.postdata;

import okhttp3.FormBody;
import org.example.e2eb.emulator.request.BasePostData;

public class CompleteData extends BasePostData {

    public CompleteData(){
        data = new FormBody.Builder()
                .add("payment_method_nonce", "")
                .add("csrfToken","haha")
                .build();
    }
}
