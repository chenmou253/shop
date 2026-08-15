import { useEffect, useState } from "react";
import { Play } from "lucide-react";
import { api } from "../api";
import { PageHero } from "../components/Common";
import type { Video as VideoData } from "../types";

export function VideoPage() {
  const [items, setItems] = useState<VideoData[]>([]);
  useEffect(() => { api<{ videos: VideoData[] }>("/videos").then((data) => setItems(data.videos)); }, []);
  return <><PageHero title="Video" /><section className="shell page-space"><div className="video-grid">{items.map((item) => <a className="video-card" href={item.video_url} target="_blank" rel="noreferrer" key={item.id}><div><img src={item.cover_url} alt={item.title} /><span><Play fill="currentColor" /></span></div><h2>{item.title}</h2><p>{item.description}</p></a>)}</div></section></>;
}
