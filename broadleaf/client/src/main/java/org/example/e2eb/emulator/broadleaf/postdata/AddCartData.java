package org.example.e2eb.emulator.broadleaf.postdata;
import okhttp3.FormBody;
import org.example.e2eb.emulator.request.BasePostData;

public class AddCartData extends BasePostData{

    public AddCartData(int productId, int quantity){
        data = new FormBody.Builder()
                .add("productId", String.valueOf(productId))
                .add("quantity", String.valueOf(quantity))
                .add("csrfToken","haha")
                .build();
    }
}
