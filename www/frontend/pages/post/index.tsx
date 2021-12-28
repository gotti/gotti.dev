import {FC} from "react";
import Link from 'next/link'
import {PageHead} from "../../components/PageHead"
import {postData, fetchPosts} from "../../libs/posts"

interface Props {
  posts: postData[];
}

const Post: FC<Props> = ({posts}) => {
  return (
    <>
      <PageHead title="Post" />
      <h2>Post</h2>
      <ul>
        {posts.map(post => (
          <li key={post.path}><Link key={post.path} href={post.path}>{post.title}</Link></li>
        ))}
      </ul>
    </>
  );
};

export const getStaticProps = async () => {
  const posts = await fetchPosts();
  return {
    props: {
      posts: posts
    },
  };
};

export default Post;
