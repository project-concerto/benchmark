package org.example.e2eb.utils;


/**
 * Class used for exit the whole program anywhere
 */
public class Panic {
    /**
     * Exit program with status 1
     * @param message reasons of exit
     */
    public static void quit(String message){
        System.out.println("Quit because " + message);
        System.exit(1);
    }
}
