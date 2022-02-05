package org.example.e2eb.emulator.broadleaf.txns;

import okhttp3.OkHttpClient;
import okhttp3.Response;
import org.example.e2eb.Config;
import org.example.e2eb.emulator.BaseEmulator;
import org.example.e2eb.emulator.broadleaf.postdata.*;
import org.example.e2eb.emulator.request.RequestUtils;
import org.example.e2eb.emulator.request.RequestFailException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;

public class BLCTxn {

    private static final Logger logger = LoggerFactory.getLogger(BLCTxn.class);

    public static final String BASE_URL = Config.getOptions().protocol + Config.getOptions().host;

    public static final boolean debugEnable = Config.getOptions().debug;;

    public static boolean doTxn(BaseEmulator emulator) {
        try{
            // addCart(emulator.getOkHttpClient(), emulator.getId());
            // Current branch is used to simulate high contention
            // Use the above line of code to simulate no contention 
            addCart(emulator.getOkHttpClient(), 666);
            checkout(emulator.getOkHttpClient(), emulator.getId() + "@qq.com");
        } catch (IOException | RequestFailException e) {
            if(debugEnable){
                logger.error(e.toString());
            }
            return false;
        }
        return true;
    }

    public static void register(OkHttpClient okHttpClient, String username, String passwd) throws IOException {
        String url = BASE_URL + "/register";
        try (Response response = RequestUtils.sendPostRequest(okHttpClient, url, new RegisterData(username, passwd), null)) {
            if (debugEnable) {
                logger.info(response.toString());
            }
        }
    }

    public static void login(OkHttpClient okHttpClient, String username, String passwd) throws IOException {
        String url = BASE_URL + "/login_post.htm";
        try (Response response = RequestUtils.sendPostRequest(okHttpClient, url, new LoginData(username, passwd), null)){
            if(debugEnable){
                logger.info(response.toString());
            }
        }
    }

    public static void addCart(OkHttpClient okHttpClient, int productId) throws IOException {
        String url = BASE_URL + "/cart/add";
        try (Response response = RequestUtils.sendPostRequest(okHttpClient, url, new AddCartData(productId, 1), null)) {
            if(debugEnable){
                logger.info(response.toString());
            }
        }
    }

    public static void checkout(OkHttpClient okHttpClient, String username) throws IOException, RequestFailException {
        String url = BASE_URL + "/checkout/singleship";
        try (Response response = RequestUtils.sendPostRequest(okHttpClient, url, new SingleShipData(), null)){
            if(debugEnable){
                logger.info(response.toString());
            }
        }


        url = BASE_URL + "/checkout/payment";
        try (Response response = RequestUtils.sendPostRequest(okHttpClient, url, new PaymentData(username), null)) {
            if(debugEnable){
                logger.info(response.toString());
            }
        }


        url = BASE_URL + "/checkout/complete";
        try (Response response = RequestUtils.sendPostRequest(okHttpClient, url, new CompleteData(), null)) {
            if(debugEnable){
                logger.info(response.toString());
            }
            validCheckoutResponse(response);
        }
    }

    public static void validCheckoutResponse(Response response) throws RequestFailException {
        if(!response.request().url().toString().contains("confirmation")){
            throw new RequestFailException(response.toString());
        }
    }
}
