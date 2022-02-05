package org.example.e2eb.utils;

import java.util.HashMap;
import java.util.Random;
import org.apache.commons.math3.distribution.ZipfDistribution;
import org.example.e2eb.Config;

/**
 * Class with some useful tools
 */
public class Utils {

    public final static Random random = new Random(1314);
    public final static ZipfDistribution zipfDistribution = new ZipfDistribution(10000,
            Config.getOptions().exponent);

    /**
     * generate a random boolean value
     * @return a random boolean
     */
    public static boolean randomBoolean(){
        return random.nextBoolean();
    }

    /**
     * generate a random int variable in [low, up]
     * @param low low bound
     * @param up up bound
     * @return a random number
     */
    public static int randomInt(int low, int up){
        return low + random.nextInt(up - low + 1);
    }

    /**
     * generate a random string contains visible ascii alphabet
     * @param len length of string
     * @return generated string
     */
    public static String randomString(int len){
        StringBuilder stringBuilder = new StringBuilder();
        for(int i=0; i<len; i++){
            if(randomBoolean()){
                stringBuilder.append((char)randomInt(65, 90));
            }
            else{
                stringBuilder.append((char)randomInt(97, 122));
            }
        }
        return stringBuilder.toString();
    }

    private final static String[] names = {"We ", "I ", "They ", "He ", "She ", "Jack ", "Jim ", "Bxb "};
    private final static String[] verbs = {"was ", "is ", "were ", "are ", "do ", "does ", "doing ", "done ", "did "};
    private final static String[] nouns = {"playing a game ", "watching television ", "talking ", "dancing ",
            "speaking ", "like ", "love ", "fuck "};

    public static String randomSentence(int len){
        StringBuilder stringBuilder = new StringBuilder();
        while(true){
            stringBuilder.append(names[randomInt(0, names.length - 1)]);
            stringBuilder.append(verbs[randomInt(0, verbs.length - 1)]);
            stringBuilder.append(nouns[randomInt(0, nouns.length - 1)]);
            if(stringBuilder.length() >= len){
                return stringBuilder.substring(0, len);
            }
        }

    }

    private static final HashMap<Class, String> simpleClassName = new HashMap<>();

    public static String getSimpleClassName(Class aClass){
        String simpleName = simpleClassName.get(aClass);
        if(simpleName == null){
            String fullName = aClass.getName();
            simpleName = fullName.substring(fullName.lastIndexOf(".") + 1);
            simpleClassName.put(aClass, simpleName);
        }
        return simpleName;
    }
}
