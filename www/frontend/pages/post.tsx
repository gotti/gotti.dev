import { FC } from "react";
import * as yaml from "js-yaml"
import Link from 'next/link'
import matter from "gray-matter"
import { GetStaticProps } from "next";

interface Post {
  title: string;
  url: string;
}

interface Props {
  posts: Post[];
}

const Post: FC<Props> = ({ posts }) => {
    return (
        <>
          <h2>Post</h2>
          <ul>
          {posts.map(post => (
            <li><Link key={post.title} href={post.url}>{post.title}</Link></li>
          ))}
          </ul>
        </>
    );
};

export const getStaticProps = async () => {
  const res = await fetch("https://raw.githubusercontent.com/gotti/blog/main/contents/blog.yaml").then(res => res.blob()).then(blob => blob.text())
  const y = yaml.load(res)["posts"]
  const paths = y.map(post => {
      return post.slice(1)
    }
   )
  const posts_b = await paths.map(async p => {
    const res = await fetch("https://raw.githubusercontent.com/gotti/blog/main/contents"+p+"/index.md").then(res => res.blob()).then(blob => blob.text())
    const y = matter(res)
    const ret = {
        title: y.data["title"],
        url: p,
    }
    return ret
  })
  const posts = await Promise.all(posts_b)
  console.log(posts)
  return {
      props: {
        posts: posts
        },
    };
};

export default Post;
