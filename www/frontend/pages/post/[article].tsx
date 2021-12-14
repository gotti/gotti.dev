import { GetStaticPaths, GetStaticProps } from 'next'
import * as yaml from "js-yaml"
import matter from "gray-matter"
import { marked } from "marked"

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
  page: string;
}

export const getStaticProps = async ({ params }) => {
  console.log("params",params)
  const a = await fetch("https://raw.githubusercontent.com/gotti/gotti.dev/main/www/contents/post/"+params.article+"/index.md").then(res => res.blob()).then(blob => blob.text())
  const b = matter(a)
  return {
    props: { page: b.content }
  }
}

const Article: NextPage<Props> = ({ page }) => {
  return (
    <>
    <div className="postBody">
      <div dangerouslySetInnerHTML={{__html: marked(page)}}></div>
    </div>
    </>
  )
}

export default Article
