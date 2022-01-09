import { GetStaticPaths, GetStaticProps } from "next";
import { useEffect, useState } from "react";
import { marked } from "marked";
import { PageHead } from "../../components/PageHead";
import { TwitterShareButton, TwitterIcon } from "react-share";
import hljs from "highlight.js";
import cheerio from "cheerio";
import { postData, Tag, fetchPosts, getTags } from "../../libs/posts";
import { BiPurchaseTagAlt } from "react-icons/bi";
import { BlogSummary } from "../../components/BlogSummary";
import "highlight.js/styles/github-dark-dimmed.css";

export const getStaticPaths: GetStaticPaths = async () => {
  const posts = await fetchPosts();
  const tags = getTags(posts);
  const paths = tags.tags.map((t) => {
    return `/tags/${t.name}`;
  });
  return {
    paths,
    fallback: false,
  };
};

interface Props {
  tag: Tag;
}

export const getStaticProps = async ({ params }) => {
  const urltag = params.tag;
  const posts = await fetchPosts();
  const tags = getTags(posts);
  let tag: Tag;
  for (const t of tags.tags) {
    if (t.name == urltag) {
      tag = t;
    }
  }
  return {
    props: { tag },
  };
};

const Tag: NextPage<Props> = ({ tag }) => {
  return (
    <>
      <PageHead title={tag.name} />
      <h1>
        <BiPurchaseTagAlt />
        {tag.name}
      </h1>
      <BlogSummary posts={tag.posts} />
    </>
  );
};

export default Tag;
