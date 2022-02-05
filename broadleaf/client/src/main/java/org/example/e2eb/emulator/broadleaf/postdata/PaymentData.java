package org.example.e2eb.emulator.broadleaf.postdata;

import okhttp3.FormBody;
import org.example.e2eb.emulator.request.BasePostData;

public class PaymentData extends BasePostData {
    
    public PaymentData(String email) {
        data = new FormBody.Builder()
                .add("paymentToken"," 4111111111111111|Hotsauce Connoisseur|01/99|123")
                .add("customerPaymentId","")
                .add("shouldUseCustomerPayment", "false")
                .add("shouldSaveNewPayment", "true")
                .add("_shouldUseShippingAddress", "on")
                .add("address.isoCountryAlpha2", "US")
                .add("address.fullName", "" )
                .add("address.addressLine1", "" )
                .add("address.addressLine2", "" )
                .add("address.city", "" )
                .add("address.stateProvinceRegion", "" )
                .add("address.postalCode", "" )
                .add("address.phonePrimary", "" )
                .add("emailAddress", email)
                .add("csrfToken","haha")
                .build();
    }
}
