import { useEffect, useState } from "react";
import { ChevronLeft, ChevronRight, Mail, MessageCircle, Plus } from "lucide-react";
import { Link, useParams } from "react-router-dom";
import { api } from "../api";
import { CategorySidebar, PageHero, ProductCard, SectionTitle } from "../components/Common";
import { useInquiry } from "../context/InquiryContext";
import { useSite } from "../context/SiteContext";
import type { Product } from "../types";

type Response = { product: Product; related: Product[]; previous?: Product; next?: Product };

export function ProductDetail() {
  const { slug } = useParams();
  const [data, setData] = useState<Response | null>(null);
  const [activeImage, setActiveImage] = useState("");
  const { add } = useInquiry();
  const { settings } = useSite();
  useEffect(() => { api<Response>(`/products/${slug}`).then((row) => { setData(row); setActiveImage(row.product.main_image); }); }, [slug]);
  if (!data) return <div className="state-page"><span className="loader" />Loading...</div>;
  const { product } = data;
  const images = [{ id: 0, image_url: product.main_image, alt_text: product.name }, ...(product.images || [])].filter((item, index, list) => list.findIndex((x) => x.image_url === item.image_url) === index);
  const specs = [["Material", product.material], ["Size", product.size], ["Thread", product.thread_standard], ["Pressure", product.pressure_rating], ["Temperature", product.temperature_range], ["MOQ", product.moq], ["Standard", product.standard], ["Application", product.application]].filter(([, v]) => v);
  const whatsappNumber = settings.whatsapp?.replace(/\D/g, "") || "";
  return <><PageHero title={product.name} crumbs={["Products", product.category?.name || ""]} /><div className="shell two-column page-space"><section className="main-column"><div className="product-detail-top"><div className="gallery"><div className="gallery-main"><img src={activeImage} alt={product.name} /></div><div className="thumbnails">{images.map((image) => <button className={activeImage === image.image_url ? "active" : ""} key={image.id} onClick={() => setActiveImage(image.image_url)}><img src={image.image_url} alt={image.alt_text} /></button>)}</div></div><div className="product-summary"><span className="eyebrow">{product.category?.name}</span><h1>{product.name}</h1><p>{product.summary}</p><div className="quick-specs">{specs.slice(0, 5).map(([label, value]) => <div key={label}><b>{label}</b><span>{value}</span></div>)}</div><button className="button large" onClick={() => add(product)}><Plus size={19} /> Add to Inquiry</button><Link className="button outline large" to={`/feedback?product=${product.id}`}><Mail size={19} /> Send Inquiry</Link>{whatsappNumber && <a className="whatsapp-link" href={`https://wa.me/${whatsappNumber}`}><MessageCircle size={18} /> Discuss this product on WhatsApp</a>}</div></div><div className="detail-section"><SectionTitle>Description</SectionTitle>{specs.length > 0 && <table className="spec-table"><tbody>{specs.map(([label, value]) => <tr key={label}><th>{label}</th><td>{value}</td></tr>)}</tbody></table>}<div className="rich-content" dangerouslySetInnerHTML={{ __html: product.description }} /><div className="hot-tags"><b>Hot Tags:</b>{product.hot_tags.split(",").filter(Boolean).map((tag) => <Link key={tag} to={`/search?q=${encodeURIComponent(tag.trim())}`}>{tag.trim()}</Link>)}</div></div><nav className="product-prev-next"><span>{data.previous?.slug && <Link to={`/product/${data.previous.slug}`}><ChevronLeft /> {data.previous.name}</Link>}</span><span>{data.next?.slug && <Link to={`/product/${data.next.slug}`}>{data.next.name} <ChevronRight /></Link>}</span></nav>{data.related.length > 0 && <section className="related"><SectionTitle>You Might Also Like</SectionTitle><div className="product-grid related-grid">{data.related.map((p) => <ProductCard product={p} key={p.id} />)}</div></section>}</section><CategorySidebar latest={data.related} /></div></>;
}
