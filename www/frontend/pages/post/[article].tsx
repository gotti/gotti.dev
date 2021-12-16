import { GetStaticPaths, GetStaticProps } from 'next'
import * as yaml from "js-yaml"
import matter from "gray-matter"
import { marked } from "marked"
import { PageHead } from "../../components/PageHead"
import { TwitterShareButton, TwitterIcon } from 'react-share';

export const getStaticPaths: GetStaticPaths = async () => {
  const res = await fetch("https://raw.githubusercontent.com/gotti/gotti.dev/main/www/contents/blog.yaml").then(res => res.blob()).then(blob => blob.text())
  console.log(res)
  const y = yaml.load(res)["posts"]
  const paths = y.map(post => {
      console.log(post)
      return post.slice(1)
    }
   )
  return {
    paths,
    fallback: false
  }
}

interface Props {
  title: string;
  page: string;
  url: string;
}

export const getStaticProps = async ({ params }) => {
  console.log("params",params)
  const a = await fetch("https://raw.githubusercontent.com/gotti/gotti.dev/main/www/contents/post/"+params.article+"/index.md").then(res => res.blob()).then(blob => blob.text())
  const b = matter(a)
  const url = "https://gotti.dev/post/"+params.article;
  return {
    props: { title: b.data["title"], page: b.content, url: url}
  }
}

const Article: NextPage<Props> = ({title, page, url}) => {
  return (
    <>
    <PageHead title={title}/>
    <TwitterShareButton title={title} url={url}>
      <TwitterIcon/>
    </TwitterShareButton>
    <div className="postBody">
      <div dangerouslySetInnerHTML={{__html: marked(page)}}></div>
    </div>
    </>
  )
}

export default Article
