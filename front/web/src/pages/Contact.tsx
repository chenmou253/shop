import { Building2, Mail, MapPin, MessageCircle, Phone } from "lucide-react";
import { PageHero } from "../components/Common";
import { useSite } from "../context/SiteContext";

export function Contact() {
  const { settings } = useSite();
  const siteName = settings.site_name || "";
  const hasContactImage = Boolean(settings.contact_image_url);
  const hasTaiwan = Boolean(settings.taiwan_company || settings.taiwan_address || settings.taiwan_email || settings.taiwan_phone);
  return <><PageHero title="Contact Us" /><section className={`shell contact-page page-space ${hasContactImage ? "" : "no-visual"}`}>{hasContactImage && <div className="contact-visual"><img src={settings.contact_image_url} alt={siteName} /></div>}<div className="contact-details"><span className="eyebrow">We respond within 24 hours</span>{siteName && <h1>{siteName}</h1>}{settings.phone && <ContactLine icon={<Phone />} title="Tel" value={settings.phone} />}{settings.whatsapp && <ContactLine icon={<MessageCircle />} title="Phone / WhatsApp" value={settings.whatsapp} />}{settings.email && <ContactLine icon={<Mail />} title="Email" value={settings.email} />}{settings.address && <ContactLine icon={<MapPin />} title="Address" value={settings.address} />}{hasTaiwan && <div className="taiwan-card"><h2><Building2 /> Taiwan Branch Factory</h2>{settings.taiwan_company && <p><b>{settings.taiwan_company}</b></p>}{settings.taiwan_address && <p>{settings.taiwan_address}</p>}{settings.taiwan_email && <p>{settings.taiwan_email}</p>}{settings.taiwan_phone && <p>{settings.taiwan_phone}</p>}</div>}</div></section></>;
}

function ContactLine({ icon, title, value }: { icon: React.ReactNode; title: string; value: string }) { return <div className="contact-line"><span>{icon}</span><div><b>{title}</b><p>{value}</p></div></div>; }
