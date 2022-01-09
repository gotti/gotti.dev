import { FC } from "react";
import Link from "next/link";
import { PageHead } from "../../components/PageHead";
import { postData, fetchPosts } from "../../libs/posts";
import { BlogSummary } from "../../components/BlogSummary";

interface Props {
  posts: postData[];
}

const Post: FC<Props> = ({ posts }) => {
  return (
    <>
      <PageHead title="Post" imageUrl="" />
      <h2>Post</h2>
      <BlogSummary posts={posts} />
    </>
  );
};

export const getStaticProps = async () => {
  const posts = await fetchPosts();
  return {
    props: {
      posts: posts,
    },
  };
};

export default Post;
