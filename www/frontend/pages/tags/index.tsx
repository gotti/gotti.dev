import { FC } from "react";
import Link from "next/link";
import { PageHead } from "../../components/PageHead";
import { postData, fetchPosts, getTags } from "../../libs/posts";

interface Props {
  posts: postData[];
}

const Tags: FC<Props> = ({ posts }) => {
  return (
    <>
      <PageHead title="Tags" />
      <h2>Tags</h2>
      <ul>
        {getTags(posts).tags.map((t) => (
          <li key={t.name}>
            <Link href={`/tags/${t.name}`}>{t.name}</Link>
          </li>
        ))}
      </ul>
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

export default Tags;
