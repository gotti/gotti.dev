import Link from "next/link";
import { FC } from "react";
import { postData, fetchPosts } from "../../libs/posts";
import { BlogTags } from "../BlogTags";

import styles from "./summary.module.scss";

interface Props {
  posts: postData[];
}

export const BlogSummary: FC<Props> = ({ posts }) => {
  return (
    <>
      <ul className={styles.posts}>
        {posts.map((post) => (
          <li className={styles.summary} key={post.path}>
            <div className={styles.title}>
              <Link key={post.path} href={post.path}>
                {post.title}
              </Link>
            </div>
            <div>{post.date}</div>
            <BlogTags tags={post.tags} />
          </li>
        ))}
      </ul>
    </>
  );
};
