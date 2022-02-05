import Axios from "axios";
import FormData from "form-data";
import fs from "fs";

const args = process.argv.slice(2)
const http = Axios.create({
  baseURL: "http://localhost:3000",
  headers: {
    "Api-Key": "1f863193afbbe719ac251c29ee73b749a562c4395d37e5c7595986d72e4c87e6",
    "Api-Username": "zxd",
    "Content-Type": "application/json",
    Accept: "application/json",
  },
});
http.interceptors.request.use((config) => {
  if (config.data instanceof FormData) {
    Object.assign(config.headers, config.data.getHeaders());
  }
  return config;
});

const filename = "./images/" + args[0] + ".png";
const file = fs.readFileSync(filename);
const form = new FormData();
form.append("files[]", file, {
  filename,
});

http
  .post("/uploads.json", form, {
    params: {
      type: "composer",
      synchronous: true,
    },
  })
  .then(({ data }) => {
    console.info("\""+args[0]+"\"" + ":" +  JSON.stringify(data, null, 2)+",");
    return {
      url: data.url,
    };
  })
  .catch((e) => {
    console.error(
      "Error uploading file to Discourse",
      JSON.stringify(e, null, 2)
    );
    throw e;
  });