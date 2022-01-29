import { FC } from "react";
import Link from "next/link";
import { PageHead } from "../components/PageHead";
import { OtakuItem } from "../components/OtakuItem"

const Page: FC = () => {
  return (
    <>
      <PageHead title="Env" imageUrl="" />
      <h1>Env</h1>
      ここは環境を自慢するページです。
      <OtakuItem name="全体" imagePath="/images/environment/overall.webp"
        description="机が汚いのは気にしないで。"/>
      <OtakuItem name="LG 43インチモニタ" imagePath="/images/environment/monitor-LG-43inch.webp"
        description="LGのクソデカモニタ。FHDの21.5インチを4枚並べるのと同等"/>
      <OtakuItem name="Acerのモニタ" imagePath="/images/environment/monitor-acer.webp"
        description="安物"/>
      <OtakuItem name="キーボード" imagePath="/images/environment/keyboard.webp"
        description="NiZのキーボードとジャンクトラックボール。ThinkPadみたいな感じにデスクトップを触りたかったのでキーボードの下にトラックボール置いてる。マウス触るのはゲームするときだけ。"/>
      <OtakuItem name="マイク" imagePath="/images/environment/mic.webp"
        description="中華の謎マイク、2千円ぐらい。" />
      <OtakuItem name="quince" imagePath="/images/environment/quince.webp"
        description="メイン機">
        <ul>
        <li>OS: Arch Linux</li>
        <li>お迎えは2020年ぐらい</li>
        </ul>
      </OtakuItem>
      <OtakuItem name="camphor" imagePath="/images/environment/camphor.webp"
        description="メインサーバ">
        <ul>
        <li>OS: Arch Linux</li>
        <li>ファイル鯖や録画鯖</li>
        <li>お迎えは2017年ぐらい</li>
        </ul>
      </OtakuItem>
      <OtakuItem name="olive" imagePath="/images/environment/olive.webp"
        description="VMサーバ">
        <ul>
        <li>OS: proxmox</li>
        <li>VMいっぱい生やしくん</li>
        <li>お迎えは2022年</li>
        </ul>
      </OtakuItem>
      <OtakuItem name="larch" imagePath="/images/environment/larch.webp">
        <ul>
        <li>OS: win10(おうち唯一)</li>
        <li>VMいっぱい生やしくん</li>
        <li>お迎えは2020年</li>
        </ul>
      </OtakuItem>
      <OtakuItem name="muffin" imagePath="/images/environment/muffin.webp"
        description="ThinkPad X390">
        <ul>
        <li>OS: Arch Linux / win10</li>
        <li>お迎えは2020年</li>
        </ul>
      </OtakuItem>
      <OtakuItem name="ルータ" imagePath="/images/environment/network.webp"
        description="上はmikrotikの3万の10Gスイッチ(ルータ機能つき)、下はRTX1200。
        quince、camphor、oliveの3台は10G喋れる。"/>
      <OtakuItem name="プリンター" imagePath="/images/environment/printer.webp"
        description="brotherのレーザープリンター">
        <ul>
        <li>安い、印刷速い、トナーはめちゃ持つ</li>
        <li>Linuxドライバ配ってる(最高)</li>
        <li>白黒しか印刷できない</li>
        </ul>
      </OtakuItem>
      <OtakuItem name="AP" imagePath="/images/environment/AP.webp"
        description="UniFi nanoHD">
        <ul>
        <li>もらった</li>
        <li>VLAN喋れて無限にSSID生やせるので良い</li>
        <li>これはwifi6非対応だけど対応してるモデルもあるらしい</li>
        </ul>
      </OtakuItem>
    </>
  );
};

export default Page;
