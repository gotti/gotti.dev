import styles from "./tags.module.scss"
import Link from 'next/link'
import {FC} from "react";
import {BiPurchaseTagAlt} from 'react-icons/bi';

interface Props {
  tags: string[];
}

export const BlogTags: FC<Props> = ({tags}) => {
  return (
    <>
      <div className={styles.tags}>
        <BiPurchaseTagAlt />
        {tags.map(tag => (
          <a className={styles.tag} href={`/tags/${tag}`} key={tag}>{tag}</a>
        ))}
      </div>
    </>
  );
}
