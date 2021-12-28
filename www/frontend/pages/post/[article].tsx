import {GetStaticPaths, GetStaticProps} from 'next';
import {useEffect, useState} from 'react';
import {marked} from "marked";
import {PageHead} from "../../components/PageHead";
import {TwitterShareButton, TwitterIcon} from 'react-share';
import {postData, fetchPost, fetchPathList} from "../../libs/posts";
import hljs from "highlight.js";
import 'highlight.js/styles/github.css';

export const getStaticPaths: GetStaticPaths = async () => {
  const posts = await fetchPathList();
  const paths = posts.map((p) => {return `/post/${p}`;})
  console.log("posts");
  console.log(paths);
  return {
    paths,
    fallback: false
  }
}

const renderMD = (text: string): string => {
  marked.setOptions({
    langPrefix: "",
    highlight: (code: string, lang: string) => {
      return hljs.highlightAuto(code, [lang]).value
    }
  });
  return marked(text);
};

interface Props {
  post: postData;
}

export const getStaticProps = async ({params}) => {
  console.log("params", params.article);
  const post = await fetchPost(params.article);
  console.log(post);
  return {
    props: {post}
  }
}

const Article: NextPage<Props> = ({post}) => {
  return (
    <>
      <PageHead title={post.title} />
      <div className="postBody">
        <div dangerouslySetInnerHTML={{__html: renderMD(post.text)}}></div>
      </div>
      <TwitterShareButton title={post.title} url={post.url}>
        <TwitterIcon size={32} round={true} />
      </TwitterShareButton>
    </>
  )
}

export default Article
