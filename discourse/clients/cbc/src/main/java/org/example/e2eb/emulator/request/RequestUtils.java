package org.example.e2eb.emulator.request;

import okhttp3.*;
import org.example.e2eb.utils.Panic;

import javax.net.ssl.*;
import java.io.IOException;
import java.security.KeyManagementException;
import java.security.NoSuchAlgorithmException;
import java.security.SecureRandom;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;


/**
 * Class used to send/receive request/response
 */
public class RequestUtils {

    private static Headers getDefaultHeader() {
        return new Headers.Builder()
                .add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36")
                .add("Connection", "keep-alive")
                .build();
    }

    /**
     * Send a get request
     *
     * @param okHttpClient client user
     * @param url          url used to send request
     * @param headers      headers, use default headers if null
     * @return the response of request
     * @throws IOException
     */
    public static Response sendGetRequest(OkHttpClient okHttpClient, String url, Headers headers) throws IOException {
        if (headers == null) {
            headers = getDefaultHeader();
        }
        Request request = new Request.Builder()
                .headers(headers)
                .url(url)
                .get()
                .build();
        Call call = okHttpClient.newCall(request);
        return call.execute();
    }

    /**
     * Send a post request
     *
     * @param okHttpClient client used
     * @param url          url used to send request
     * @param postData     post data
     * @param headers      headers, use default headers if null
     * @return the response of request
     * @throws IOException
     */
    public static Response sendPostRequest(OkHttpClient okHttpClient, String url, BasePostData postData, Headers headers) throws IOException {
        if (headers == null) {
            headers = getDefaultHeader();
        }
        Request request = new Request.Builder()
                .headers(headers)
                .url(url)
                .post(postData.getData())
                .build();
        Call call = okHttpClient.newCall(request);
        return call.execute();
    }

    /**
     * Send a post request
     *
     * @param okHttpClient client used
     * @param url          url used to send request
     * @param postData     put data
     * @param headers      headers, use default headers if null
     * @return the response of request
     * @throws IOException
     */
    public static Response sendPutRequest(OkHttpClient okHttpClient, String url, BasePostData postData, Headers headers) throws IOException {
        if (headers == null) {
            headers = getDefaultHeader();
        }
        Request request = new Request.Builder()
                .headers(headers)
                .url(url)
                .put(postData.getData())
                .build();
        Call call = okHttpClient.newCall(request);
        return call.execute();
    }

    private static OkHttpClient.Builder builder = null;

    public static OkHttpClient.Builder getInsecureOkHttpClientBuilder() {
        if(builder == null){
            TrustManager[] trustAllCerts = new TrustManager[]{trustManager()};

            SSLContext sslContext = null;
            try {
                sslContext = SSLContext.getInstance("SSL");
                sslContext.init(null, trustAllCerts, new SecureRandom());
            } catch (NoSuchAlgorithmException | KeyManagementException e) {
                Panic.quit(e.toString());
            }

            assert sslContext != null;
            SSLSocketFactory sslSocketFactory = sslContext.getSocketFactory();

            Dispatcher dispatcher = new Dispatcher();
            dispatcher.setMaxRequests(10240);
            dispatcher.setMaxRequestsPerHost(10240);

            builder = new OkHttpClient.Builder().dispatcher(dispatcher);
            builder.sslSocketFactory(sslSocketFactory, (X509TrustManager) trustAllCerts[0]);
            builder.hostnameVerifier(hostnameVerifier());
        }
        return builder;
    }

    private static TrustManager trustManager() {
        return new X509TrustManager() {
            @Override
            public void checkClientTrusted(X509Certificate[] chain, String authType) throws CertificateException {
            }

            @Override
            public void checkServerTrusted(X509Certificate[] chain, String authType) throws CertificateException {
            }

            @Override
            public X509Certificate[] getAcceptedIssuers() {
                return new X509Certificate[]{};
            }
        };
    }

    private static HostnameVerifier hostnameVerifier() {
        return new HostnameVerifier() {
            @Override
            public boolean verify(String hostname, SSLSession session) {
                return true;
            }
        };
    }
}
