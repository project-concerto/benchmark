package org.example.e2eb.emulator.request;

public class RequestFailException extends Exception {
    public RequestFailException(String message){
        super(message);
    }
}