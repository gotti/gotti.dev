import {GetStaticPaths, GetStaticProps} from 'next'
import * as yaml from "js-yaml"
import matter from "gray-matter"
import {marked} from "marked"
import {PageHead} from "../../components/PageHead"
import {TwitterShareButton, TwitterIcon} from 'react-share';
import {postData, fetchPost, fetchPathList} from "../../libs/posts"
import {buildPostURL} from '../../libs/settings'

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
        <div dangerouslySetInnerHTML={{__html: marked(post.text)}}></div>
      </div>
      <TwitterShareButton title={post.title} url={post.url}>
        <TwitterIcon size={32} round={true} />
      </TwitterShareButton>
    </>
  )
}

export default Article
