import Link from "next/link";

export const Menu: NextPage = () => {
  return (
    <>
      <ul>
        <li>
          <Link href="/">Home</Link>
        </li>
        <li>
          <Link href="/post">Post</Link>
        </li>
        <li>
          <Link href="/environment">Env</Link>
        </li>
      </ul>
    </>
  );
};
