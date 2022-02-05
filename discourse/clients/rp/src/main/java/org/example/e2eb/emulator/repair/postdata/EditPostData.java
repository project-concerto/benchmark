package org.example.e2eb.emulator.repair.postdata;

import com.fasterxml.jackson.databind.ObjectMapper;
import okhttp3.FormBody;
import org.example.e2eb.emulator.request.BasePostData;
import org.example.e2eb.utils.Utils;

import java.io.IOException;
import java.nio.file.Paths;
import java.util.Map;

public class EditPostData extends BasePostData {
    private static Map<String, Map<String, String>> images;
    static {
        ObjectMapper objectMapper = new ObjectMapper();
        try {
            images = objectMapper.readValue(Paths.get("./images.json").
                    toFile(), Map.class);
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
    public EditPostData(int eId) {
        String raw = Utils.randomSentence(50);
        String imageId = String.valueOf((eId+7)/8);
        raw += String.format("![%s|%sx%s](%s)", imageId,
                images.get(imageId).get("thumbnail_width"),
                images.get(imageId).get("thumbnail_height"),
                images.get(imageId).get("short_url"));
        String payload = String.format("{\n" +
                "\"raw\": \"%s\",\n" +
                "\"edit_reason\": \"%s\"\n" +
                "}", raw, Utils.randomSentence(10));
        data = new FormBody.Builder()
                .add("post[raw]", raw)
                .build();
    }
}
