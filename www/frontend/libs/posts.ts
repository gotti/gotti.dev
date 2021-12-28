import {buildPostListURL, buildPostURL, buildSiteURL} from "./settings"
import * as yaml from "js-yaml"
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

export const fetchPost = async (post: string): Promise<postData> => {
  console.log(post);
  console.log(buildPostURL(post));
  const p = await fetch(buildPostURL(post));
  const rawpost = await p.text();
  const mpost = matter(rawpost);
  const ret: postData = {
    title: mpost.data["title"],
    date: mpost.data["date"],
    tags: mpost.data["tags"],
    text: mpost.content,
    url: buildSiteURL(post),
    name: post,
    path: `./post/${post}`,
  };
  return ret;
}

export const fetchPosts = async (): Promise<postData[]> => {
  const postlist = await fetchPathList()
  const posts = postlist.map(async (post: string) => {
    const p = await fetch(buildPostURL(post));
    const rawpost = await p.text();
    const mpost = matter(rawpost);
    const ret: postData = {
      title: mpost.data["title"],
      date: mpost.data["date"],
      tags: mpost.data["tags"],
      text: mpost.content,
      url: buildSiteURL(post),
      name: post,
      path: `./post/${post}`,
    };
    return ret;
  })
  const ret = Promise.all(posts)
  return ret;
}
