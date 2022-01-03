import Head from 'next/head'
import Link from 'next/link'
import {FC} from "react";
import {postData, fetchPosts} from "../../libs/posts"

import {BiPurchaseTagAlt} from 'react-icons/bi';

import styles from "./summary.module.scss"

interface Props {
  posts: postData[];
}

export const BlogSummary: FC<Props> = ({posts}) => {
  return (
    <>
      <ul className={styles.posts}>
        {posts.map(post => (
          <li className={styles.summary} key={post.path}>
            <div className={styles.title}><Link key={post.path} href={post.path}>{post.title}</Link></div>
            <div>{post.date}</div>
            <div className={styles.tags}>
              <BiPurchaseTagAlt />
                {post.tags.map(tag => (
                  <a className={styles.tag} href={`/tags/${tag}`}>{tag}</a>
                ))}
            </div>
          </li>
        ))}
      </ul>
    </>
  )
}

