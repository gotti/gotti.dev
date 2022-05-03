import canvas from "canvas";
import path from "path";
import fs from "fs/promises";

const ogpHeight = 630;
const ogpWidth = 1200;

const imageMaxHeight = 630;
const imageMaxWidth = 1200;

export const generateOGPWithImage = async (
  image: canvas.Image,
  filename: string,
  text: string
) => {
  const cvs = canvas.createCanvas(ogpWidth, ogpHeight);
  const ctx = cvs.getContext("2d");
  const wph = image.width / image.height;
  ctx.font = `50px sans-serief`;
  const stage1height = imageMaxHeight;
  const stage1width = wph * imageMaxHeight;
  const stage2height =
    stage1width > imageMaxWidth ? imageMaxHeight : stage1height;
  const stage2width =
    stage1width > imageMaxWidth ? wph * imageMaxHeight : stage1width;
  ctx.drawImage(
    image,
    ogpWidth / 2 - stage2width / 2,
    30,
    stage2width,
    stage2height
  );
  const l = ctx.measureText(text);
  ctx.fillText(text, ogpWidth / 2 - l.width / 2, ogpHeight - 100);

  const buf = cvs.toBuffer("image/jpeg", { quality: 0.5 });
  await fs.writeFile(`./out/${filename}.png`, buf);
};

export const generateOGPWithOutImage = async (
  filename: string,
  text: string
) => {
  const cvs = canvas.createCanvas(ogpWidth, ogpHeight);
  const ctx = cvs.getContext("2d");
  const l = ctx.measureText(text);
  ctx.fillText(text, ogpWidth / 2 - l.width / 2, ogpHeight - 100);

  const buf = cvs.toBuffer("image/jpeg", { quality: 0.5 });
  await fs.writeFile(`./out/${filename}.png`, buf);
};

export const generateOGPWithImage2 = async (
  image: canvas.Image,
  filename: string
) => {
  const cvs = canvas.createCanvas(ogpWidth, ogpHeight);
  const ctx = cvs.getContext("2d");
  const wph = image.width / image.height;
  ctx.font = `50px sans-serief`;
  const stage1height = imageMaxHeight;
  const stage1width = wph * imageMaxHeight;
  const stage2height =
    stage1width > imageMaxWidth ? imageMaxHeight : stage1height;
  const stage2width =
    stage1width > imageMaxWidth ? wph * imageMaxHeight : stage1width;
  ctx.drawImage(
    image,
    ogpWidth / 2 - stage2width / 2,
    30,
    stage2width,
    stage2height
  );
  const buf = cvs.toBuffer("image/jpeg", { quality: 0.5 });
  await fs.writeFile(`./out/${filename}.png`, buf);
};
