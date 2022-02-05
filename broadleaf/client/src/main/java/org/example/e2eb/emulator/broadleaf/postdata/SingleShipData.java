package org.example.e2eb.emulator.broadleaf.postdata;

import okhttp3.FormBody;
import org.example.e2eb.emulator.request.BasePostData;

public class SingleShipData extends BasePostData {

    public SingleShipData() {
        data = new FormBody.Builder()
                .add("address.isoCountryAlpha2", "US")
                .add("address.fullName", "zhangxiaodong")
                .add("address.addressLine1", "4301")
                .add("address.addressLine2", "x36-6031")
                .add("address.city", "shanghai")
                .add("address.stateProvinceRegion", "SK")
                .add("address.postalCode", "00001")
                .add("address.phonePrimary.phoneNumber", "13571219429")
                .add("saveAsDefault", "false")
                .add("fulfillmentOptionId", "1")
                .add("csrfToken", "haha")
                .build();
    }
}
