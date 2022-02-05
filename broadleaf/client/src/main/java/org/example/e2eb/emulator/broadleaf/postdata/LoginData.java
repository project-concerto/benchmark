package org.example.e2eb.emulator.broadleaf.postdata;

import okhttp3.FormBody;
import org.example.e2eb.emulator.request.BasePostData;

public class LoginData extends BasePostData {

    public LoginData(String username, String passwd){
        data = new FormBody.Builder()
                .add("username",username)
                .add("password",passwd)
                .add("remember-me","false")
                .add("csrfToken","haha")
                .build();
    }

}
