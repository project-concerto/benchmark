package org.example.e2eb.emulator.broadleaf.postdata;

import okhttp3.FormBody;
import org.example.e2eb.emulator.request.BasePostData;

public class RegisterData extends BasePostData {

    public RegisterData(String email, String passwd){
        data = new FormBody.Builder()
                .add("redirectUrl", "")
                .add("customer.emailAddress", email)
                .add("customer.firstName", "")
                .add("customer.lastName", "qq")
                .add("password", passwd)
                .add("passwordConfirm", passwd)
                .add("csrfToken","haha")
                .build();
    }
}
