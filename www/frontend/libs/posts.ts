import {buildPostListURL, buildPostURL, buildFileURL , buildSiteURL} from "./settings"
import * as yaml from "js-yaml"
import request from "request";
import matter from "gray-matter"

export interface posts {
  posts: postData[];
}
export interface postData {
  title: string;
  date: Date;
  tags: string[];
  text: string;
  url: string;
  name: string;
  path: string;
  ogpImagePath: string;
}

export const fetchPathList = async (): Promise<string[]> => {
  const postlist_p = await fetch(buildPostListURL());
  const postlist = await postlist_p.text();
  const posts: string[] = yaml.load(postlist)["posts"];
  const postNames = posts.map(post => {
    return post.replace("./post/", "");
  })
  return postNames;
};

const mattertoPostData = (post: string, mpost: matter.GrayMatterFile<string>, ogppath: string): postData => {
  console.log(mpost.data["tags"]);
  //ほんとに引数の型あってる？
  const ret: postData = {
    title: mpost.data["title"],
    date: mpost.data["date"],
    tags: mpost.data["tags"],
    text: mpost.content,
    url: buildSiteURL(post),
    name: post,
    path: `/post/${post}`,
    ogpImagePath: `https://gotti.dev/post/${post}/${ogppath}`,
  };
  return ret;
}

export const fetchPost = async (post: string): Promise<postData> => {
  const p = await fetch(buildPostURL(post));
  const rawpost = await p.text();
  const mpost = matter(rawpost);
  const ipath = mpost.content.match(/\!\[.+\]\((.+)\)/);
  const imagepath = ipath === null ? "" : ipath[1];
  const ret = mattertoPostData(post, mpost, imagepath);
  console.log(ret);
  return ret;
}

export const fetchPosts = async (): Promise<postData[]> => {
  const postlist = await fetchPathList()
  const posts = postlist.map(async (post: string) => {
    const ret = await fetchPost(post);
    return ret;
  })
  const ret = Promise.all(posts)
  return ret;
}

export interface Tag {
  name: string;
  posts: postData[];
}

interface Tags {
  tags: Tag[];
}

export const getTags = (posts: postData[]): Tags => {
  let tags = new Map<string, postData[]>();
  for (const p of posts) {
    for (const t of p.tags) {
      if (tags[t] == undefined) {
        tags[t] = [];
      }
      tags[t].push(p)
    }
  }
  let ret: Tags = {tags: []};
  for (const t of Object.keys(tags)) {
    const tmp: Tag = {name: t, posts: tags[t]}
    ret.tags.push(tmp)
  }
  return ret;
}
