import { GetStaticPaths, GetStaticProps } from "next";
import { marked } from "marked";
import { BlogTags } from "../../components/BlogTags";
import { PageHead } from "../../components/PageHead";
import { TwitterShareButton, TwitterIcon } from "react-share";
import { postData, fetchPost, fetchPathList } from "../../libs/posts";
import hljs from "highlight.js";
import cheerio from "cheerio";
import "highlight.js/styles/github-dark-dimmed.css";

export const getStaticPaths: GetStaticPaths = async () => {
  const posts = await fetchPathList();
  const paths = posts.map((p) => {
    return `/post/${p}`;
  });
  console.log("posts");
  console.log(paths);
  return {
    paths,
    fallback: false,
  };
};

const setInpageLink = (html: string): string => {
  const headings = cheerio.load(html);
  headings("h1, h2, h3").replaceWith((_i, elm): any => {
    const tagId = headings(elm).attr("id");
    headings(elm).wrap(`<a href="#${tagId}"></a>`);
  });
  return headings.html();
};

const renderMD = (text: string): string => {
  marked.setOptions({
    langPrefix: "",
    highlight: (code: string, lang: string) => {
      return hljs.highlightAuto(code, [lang]).value;
    },
  });
  const html: string = marked(text);
  return html;
};

interface Props {
  post: postData;
}

export const getStaticProps = async ({ params }) => {
  const post = await fetchPost(params.article);
  return {
    props: { post },
  };
};

const Article: NextPage<Props> = ({ post }) => {
  return (
    <>
      <PageHead title={post.title} imageUrl={post.ogpImagePath} />
      <div className="postDescription">
        <h1>{post.title}</h1>
        posted on {post.date}
        <BlogTags tags={post.tags} />
      </div>
      <div className="postBody">
        <div dangerouslySetInnerHTML={{ __html: renderMD(post.text) }}></div>
      </div>
      <TwitterShareButton title={post.title} url={post.url}>
        <TwitterIcon size={32} round={true} />
      </TwitterShareButton>
    </>
  );
};

export default Article;
