package org.example.e2eb.emulator.request;

import okhttp3.RequestBody;


/**
 * Base post data class, create other post data class based
 * on this, and add data according post type
 */
public class BasePostData {
    protected RequestBody data;

    public BasePostData(){}

    public BasePostData(RequestBody data){
        this.data = data;
    }

    public RequestBody getData() {
        return data;
    }
}
