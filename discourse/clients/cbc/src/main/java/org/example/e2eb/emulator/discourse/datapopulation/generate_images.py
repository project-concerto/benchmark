from PIL import Image
if __name__ == '__main__':
    for i in range(129)[1:]:
        img = Image.new("RGBA", (4000, 4000), "rgb({}, {}, {})".format(i, 50+i, 100+i))
        img.save("./images/{}.png".format(i))
