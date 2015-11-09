package main

import (
	"./introspect"
	"fmt"
)

func goodType(sig string) string {
	switch introspect.NewTypeMeta(sig).Type {
	case introspect.PrimitiveTypeId:
		return tTable.Get(sig)
	default:
		return "dbus::types::" + normalizeSignature(sig)
	}
}

func declareType(sigs []string, filedPrefix string) string {
	var r = ""
	for i, sig := range sigs {
		if i != 0 {
			r += ", "
		}
		r = r + goodType(sig)
		if filedPrefix != "" {
			r = r + fmt.Sprintf(" %s%d", filedPrefix, i)
		}
	}
	return r
}

var s_common = `
template<int index, typename T1=uchar, typename T2=uchar, typename T3=uchar, typename T4=uchar, typename T5=uchar, typename T6=uchar, typename T7=uchar, typename T8=uchar>
struct SelectBase {
    typedef void Type;
};

template<typename T1, typename T2, typename T3, typename T4, typename T5, typename T6, typename T7, typename T8>
struct SelectBase<0, T1, T2, T3, T4, T5, T6, T7, T8> {
    typedef T1 Type;
};

template<typename T1, typename T2, typename T3, typename T4, typename T5, typename T6, typename T7, typename T8>
struct SelectBase<1, T1, T2, T3, T4, T5, T6, T7, T8> {
    typedef T2 Type;
};

template<typename T1, typename T2, typename T3, typename T4, typename T5, typename T6, typename T7, typename T8>
struct SelectBase<2, T1, T2, T3, T4, T5, T6, T7, T8> {
    typedef T3 Type;
};
template<typename T1, typename T2, typename T3, typename T4, typename T5, typename T6, typename T7, typename T8>
struct SelectBase<3, T1, T2, T3, T4, T5, T6, T7, T8> {
    typedef T4 Type;
};
template<typename T1, typename T2, typename T3, typename T4, typename T5, typename T6, typename T7, typename T8>
struct SelectBase<4, T1, T2, T3, T4, T5, T6, T7, T8> {
    typedef T5 Type;
};
template<typename T1, typename T2, typename T3, typename T4, typename T5, typename T6, typename T7, typename T8>
struct SelectBase<5, T1, T2, T3, T4, T5, T6, T7, T8> {
    typedef T6 Type;
};

template<typename T1, typename T2, typename T3, typename T4, typename T5, typename T6, typename T7, typename T8>
struct SelectBase<-1, T1, T2, T3, T4, T5, T6, T7, T8> {
    typedef T7 Type;
};



typedef QVariant (*DataConverter)(QVariant);

inline QVariant NormalConverter(QVariant v) 
{ 
    return v;
}
inline QVariant PropertyConverter(QVariant v) 
{
    QDBusVariant vv = v.value<QDBusVariant>();
    return vv.variant();
}


template<typename T1=uchar, typename T2=uchar, typename T3=uchar, typename T4=uchar, typename T5=uchar, typename T6=uchar, typename T7=uchar, typename T8=uchar> class R
    : public QDBusPendingReply<> {
    template<int index>
    struct Select: SelectBase<index, T1, T2, T3, T4, T5, T6, T7, QDBusError> {
    };
private:
    void waitForFinished() {
        QDBusPendingReply::waitForFinished();
	if (!isValid() || isError()) {
            m_hasError = true;
            m_error = error();
	    return;
        }
        m_hasError = false;
        m_error = QDBusError();
    }
    QDBusError m_error;
    DataConverter m_converter;
    bool m_hasError;
public:
    R(QDBusPendingReply<> r, DataConverter c=NormalConverter):
        QDBusPendingReply(r), m_converter(c),m_hasError(false)
    {
    }

    bool hasError() {
        if (!isFinished()) {
            waitForFinished();
        }
        return m_hasError;
    }
    QDBusError Error() {
        if (!isFinished()) {
            waitForFinished();
        }
        return m_error;
    }

    template<int index>
    typename Select<index>::Type Value() {
        if (!isFinished()) {
            waitForFinished();
            if (m_hasError) {
                return typename Select<index>::Type();
            }
        }
        QList<QVariant> args = reply().arguments();
        if (args.size() <= index) {
            m_hasError = true;
            m_error = QDBusError(QDBusError::InvalidArgs, QString("can't fetch the %1th argument, because only up to %2 arguments.").arg(index).arg(args.size()));
            return typename Select<index>::Type();
        }
        QVariant r = args[index];

        return qdbus_cast<typename Select<index>::Type>(m_converter(r));
    }

    QList<QVariant> Values() {
        QList<QVariant> ret;

        if (!isFinished()) {
            waitForFinished();
            if (m_hasError) {
                return ret;
            }
        }

        QList<QVariant> args = reply().arguments();

        switch (args.size()) {
            case 8:
		    ret.push_front(QVariant::fromValue(qdbus_cast<T8>(m_converter(args[7]))));
            case 7:
		    ret.push_front(QVariant::fromValue(qdbus_cast<T7>(m_converter(args[6]))));
            case 6:
		    ret.push_front(QVariant::fromValue(qdbus_cast<T6>(m_converter(args[5]))));
            case 5:
		    ret.push_front(QVariant::fromValue(qdbus_cast<T5>(m_converter(args[4]))));
            case 4:
		    ret.push_front(QVariant::fromValue(qdbus_cast<T4>(m_converter(args[3]))));
            case 3:
		    ret.push_front(QVariant::fromValue(qdbus_cast<T3>(m_converter(args[2]))));
            case 2:
		    ret.push_front(QVariant::fromValue(qdbus_cast<T2>(m_converter(args[1]))));
            case 1:
		    ret.push_front(QVariant::fromValue(qdbus_cast<T1>(m_converter(args[0]))));
        }
        return ret;
    }

};



static QDBusConnection detectConnection(QString addr) {
    if (addr == "session") {
	    return QDBusConnection::sessionBus();
    } else if (addr == "system") {
	    return QDBusConnection::systemBus();
    } else {
            qDebug() << "W: Are you sure using '" << addr << "' as the private dbus session?";
	    return *(new QDBusConnection(addr));
    }
}

class DBusObject: public QDBusAbstractInterface {
    Q_OBJECT
public:
    DBusObject(QObject* parent, QString service, QString path, const char* interface, const QString addr);
    ~DBusObject();

    Q_SLOT void propertyChanged(const QDBusMessage& msg);
protected:
    QDBusPendingReply<> fetchProperty(const char* name);

};



inline
DBusObject::DBusObject(QObject* parent, QString service, QString path, const char* interface, const QString addr)
:QDBusAbstractInterface(service, path, interface, detectConnection(addr), parent)
{
        if (!isValid()) {
            qDebug() << "The remote dbus object may not exists, try launch it. " << lastError();
        }
	connection().connect(this->service(), this->path(), "org.freedesktop.DBus.Properties",  "PropertiesChanged",
	    "sa{sv}as", this, SLOT(propertyChanged(QDBusMessage)));
}

inline
DBusObject::~DBusObject()
{
	connection().disconnect(service(), path(), interface(),  "PropertiesChanged",
            "sa{sv}as", this, SLOT(propertyChanged(QDBusMessage)));
}

inline
QDBusPendingReply<> DBusObject::fetchProperty(const char* name)
{
    QDBusMessage msg = QDBusMessage::createMethodCall(service(), path(),
            QLatin1String("org.freedesktop.DBus.Properties"),
            QLatin1String("Get"));
    msg << interface() << QString::fromUtf8(name);

    QDBusPendingReply<> r = connection().asyncCall(msg);

    return r;
}

inline
void DBusObject::propertyChanged(const QDBusMessage& msg)
{
    QList<QVariant> arguments = msg.arguments();
    if (3 != arguments.count())
        return;

    QVariantMap changedProps = qdbus_cast<QVariantMap>(arguments.at(1).value<QDBusArgument>());
    foreach(const QString &prop, changedProps.keys()) {
        const QMetaObject* self = metaObject();
        for (int i=self->propertyOffset(); i < self->propertyCount(); ++i) {
            QMetaProperty p = self->property(i);
            if (p.name() == prop) {
                Q_EMIT p.notifySignal().invoke(this);
            }
        }
    }
}
`
