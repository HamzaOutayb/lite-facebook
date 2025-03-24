import Link from 'next/link';
import { useEffect, useState, useRef } from 'react';
import './notification.css';

const Notifications = ({ notifications = [], Err }) => {
  const [items, setItems] = useState(notifications);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const containerRef = useRef();

  useEffect(() => {
    const fetchItems = async () => {
      setLoading(true);
      try {
        const res = await fetch(`http://localhost:8080/api/GetNotification/?page=${page}`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
        });

        const newItems = await res.json();
        if (res.status === 200 && newItems.notifications) {
          setItems((prev) => [...prev, ...newItems.notifications]);
        } else {
          console.error('Error fetching notifications:', newItems);
        }
      } catch (error) {
        console.error('Error fetching notifications:', error);
      }
      setLoading(false);
    };

    fetchItems();
  }, [page]);

  useEffect(() => {
    const handleScroll = () => {
      if (containerRef.current.scrollTop <= 20 && !loading) {
        setPage((prev) => prev + 1);
      }
    };

    const container = containerRef.current;
    if (container) container.addEventListener('scroll', handleScroll);

    return () => {
      if (container) container.removeEventListener('scroll', handleScroll);
    };
  }, [loading, items]);

  const Handlefollow = async (id) => {
    try {
      await fetch('http://localhost:8080/api/deletenotification', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ id }),
      });
      setItems((prev) => prev.filter((item) => item.id !== id));
    } catch (error) {
      console.error('Error handling follow:', error);
    }
  };

  return (
    <div className="notification-wrapper" >
      <div className="notification-container" ref={containerRef}>
        {Err && <div className="notif-err">Error loading notifications. Please try again.</div>}
        {items.map((notification, index) => {
          switch (notification.type) {
            case 'follow':
              return (
                <div key={index} className="notification-div">
                  <h1>New Follower</h1>
                  <p>You received a follow from {notification.invoker}</p>
                </div>
              );
            case 'follow-request':
              return (
                <div key={index} className="notification-div">
                  <h1>Follow Request</h1>
                  <p>{notification.invoker} sent you a follow request</p>
                  <button className="accept" onClick={() => Handlefollow(notification.id)}>Accept</button>
                  <button className="reject">Reject</button>
                </div>
              );
            case 'invitation-request':
              return (
                <div key={index} className="notification-div">
                  <h1>Invitation Request</h1>
                  <p>{notification.invoker} invited you to join the group {notification.group}</p>
                  <button className="accept">Accept</button>
                  <button className="reject">Reject</button>
                </div>
              );
            case 'joine':
              return (
                <div key={index} className="notification-div">
                  <h1>Group Joining Request</h1>
                  <p>{notification.invoker} sent a join request to {notification.group}</p>
                  <button className="accept">Accept</button>
                  <button className="reject">Reject</button>
                </div>
              );
            case 'event-created':
              return (
                <div key={index} className="notification-div">
                  <Link href={`/event/${notification.eventID}`}>
                    <h1>New Event</h1>
                  </Link>
                  <p>{notification.invoker} created an event in {notification.group}</p>
                </div>
              );
            default:
              return null;
          }
        })}
        {loading && <div className="loading">Loading more notifications...</div>}
      </div>
    </div>
  );
};

export default Notifications;
