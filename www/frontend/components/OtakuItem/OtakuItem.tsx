import Link from "next/link";
import { FC } from "react";

import styles from "./item.module.scss";

interface Props {
  name: string;
  imagePath: string;
  description?: string;
}

export const OtakuItem: FC<Props> = ({ children, name, imagePath, description}) => {
  return (
    <div className={styles.otakuItem}>
      <h2 className={styles.otakuItemTitle}>{name}</h2>
      <img className={styles.otakuItemImage} src={imagePath}/>
      <div className={styles.description}>
        {description}
      </div>
      {children}
    </div>
  );
};
