import Link from 'next/link';
import { useEffect, useState, useRef } from 'react';
import './notification.css';
import { FetchApi } from '@/app/helpers';
import { useRouter } from 'next/navigation';

const Notifications = ({ notifications = [], Err }) => {
  const [items, setItems] = useState(notifications);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const containerRef = useRef();
  const router = useRouter();

  useEffect(() => {
    const fetchItems = async () => {
      setLoading(true);
      try {
        const res = await FetchApi(`/api/GetNotification/?page=${page}`, router, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
        });

        if (!res.ok) throw new Error('Failed to fetch notifications');

        const newItems = await res.json();
        if (newItems.notifications) {
          setItems((prev) => [...prev, ...newItems.notifications]);
        }
      } catch (error) {
        console.error('Error fetching notifications:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchItems();
  }, [page, router]);

  useEffect(() => {
    const handleScroll = () => {
      if (containerRef.current && containerRef.current.scrollTop <= 20 && !loading) {
        setPage((prev) => prev + 1);
      }
    };

    const container = containerRef.current;
    container?.addEventListener('scroll', handleScroll);

    return () => container?.removeEventListener('scroll', handleScroll);
  }, [loading]);

  const HandleFollow = async (id, follower, status) => {
    try {
      const followRes = await fetch('/api/follow/decision', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ follower, status }),
        credentials: 'include',
      });

      if (!followRes.ok) throw new Error('Error processing follow request');

      await fetch('/api/deletenotification', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ id }),
        credentials: 'include',
      });

      // Remove the notification from the UI after processing
      setItems((prev) => prev.filter((item) => item.id !== id));
    } catch (error) {
      console.error('Error handling follow:', error);
    }
  };

  return (
    <div className="notification-wrapper">
      <div className="notification-container" ref={containerRef}>
        {Err && <div className="notif-err">Error loading notifications. Please try again.</div>}
        {items.map((notification, index) => {
          const { id, invoker, invoker_id, invoker_name, group, eventID, type } = notification;

          return (
            <div key={index} className="notification-div">
              {type === 'follow' && (
                <>
                  <h1>A Follow</h1>
                  <p>You received a follow from {invoker_name}</p>
                </>
              )}
              {type === 'follow-request' && (
                <>
                  <h1>Follow Request</h1>
                  <p>{invoker} sent you a follow request</p>
                  <button className="accept" onClick={() => HandleFollow(id, invoker_id, 'accepted')}>Accept</button>
                  <button className="reject" onClick={() => HandleFollow(id, invoker_id, 'rejected')}>Reject</button>
                </>
              )}
              {type === 'invitation-request' && (
                <>
                  <h1>Invitation Request</h1>
                  <p>{invoker} invited you to join the group {group}</p>
                  <button className="accept">Accept</button>
                  <button className="reject">Reject</button>
                </>
              )}
              {type === 'joine' && (
                <>
                  <h1>Group Joining Request</h1>
                  <p>{invoker} sent a join request to {group}</p>
                  <button className="accept">Accept</button>
                  <button className="reject">Reject</button>
                </>
              )}
              {type === 'event-created' && (
                <>
                  <Link href={`/event/${eventID}`}>
                    <h1>New Event</h1>
                  </Link>
                  <p>{invoker} created an event in {group}</p>
                </>
              )}
            </div>
          );
        })}
        {loading && <div className="loading">Loading more notifications...</div>}
      </div>
    </div>
  );
};

export default Notifications;
